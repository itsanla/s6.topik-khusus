import amqp, { type Channel, type ChannelModel, type ConsumeMessage } from "amqplib";

import { rabbitMqConfig } from "./env";

type RabbitState = {
  connectionPromise?: Promise<ChannelModel>;
  channelPromises: Map<string, Promise<Channel>>;
};

const globalForRabbit = globalThis as typeof globalThis & {
  __rabbitState?: RabbitState;
};

const rabbitState: RabbitState =
  globalForRabbit.__rabbitState ?? {
    channelPromises: new Map<string, Promise<Channel>>(),
  };

globalForRabbit.__rabbitState = rabbitState;

const encodeVhost = (vhost: string): string => {
  if (vhost === "/") return "%2F";
  return encodeURIComponent(vhost.replace(/^\//, ""));
};

const getRabbitUrl = (): string => {
  const { host, password, port, user, vhost } = rabbitMqConfig;
  return `amqp://${encodeURIComponent(user)}:${encodeURIComponent(password)}@${host}:${port}/${encodeVhost(vhost)}`;
};

const getConnection = async (): Promise<ChannelModel> => {
  if (!rabbitState.connectionPromise) {
    rabbitState.connectionPromise = amqp.connect(getRabbitUrl());
  }

  return rabbitState.connectionPromise;
};

export const getRabbitChannel = async (queueName: string): Promise<Channel> => {
  const existing = rabbitState.channelPromises.get(queueName);
  if (existing) {
    return existing;
  }

  const channelPromise = (async () => {
    const connection = await getConnection();
    const channel = await connection.createChannel();
    await channel.assertQueue(queueName, { durable: true });
    return channel;
  })();

  rabbitState.channelPromises.set(queueName, channelPromise);
  return channelPromise;
};

export const publishToQueue = async <T>(queueName: string, payload: T): Promise<void> => {
  const channel = await getRabbitChannel(queueName);
  const serialized = Buffer.from(JSON.stringify(payload));

  const published = channel.sendToQueue(queueName, serialized, {
    contentType: "application/json",
    persistent: true,
  });

  if (!published) {
    throw new Error("Queue publish buffer is full, please retry.");
  }
};

export const consumeQueue = async (
  queueName: string,
  onMessage: (payload: unknown, raw: ConsumeMessage) => Promise<void>,
): Promise<void> => {
  const channel = await getRabbitChannel(queueName);
  await channel.prefetch(5);

  await channel.consume(queueName, async (message) => {
    if (!message) return;

    try {
      const payload = JSON.parse(message.content.toString("utf8")) as unknown;
      await onMessage(payload, message);
      channel.ack(message);
    } catch (error) {
      console.error("Email consumer failed to process message", error);
      channel.nack(message, false, false);
    }
  });
};

export const closeRabbit = async (): Promise<void> => {
  if (!rabbitState.connectionPromise) return;

  const connection = await rabbitState.connectionPromise;
  await connection.close();

  rabbitState.connectionPromise = undefined;
  rabbitState.channelPromises = new Map<string, Promise<Channel>>();
};
