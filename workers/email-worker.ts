import { loadEnvConfig } from "@next/env";

loadEnvConfig(process.cwd());

import type { EmailPayload } from "../lib/email";
import { sendEmail } from "../lib/email";
import { rabbitMqConfig } from "../lib/env";
import { closeRabbit, consumeQueue } from "../lib/rabbitmq";

const isEmailPayload = (payload: unknown): payload is EmailPayload => {
  if (!payload || typeof payload !== "object") return false;

  const candidate = payload as Record<string, unknown>;

  const hasTo = typeof candidate.to === "string" && candidate.to.length > 0;
  const hasSubject =
    typeof candidate.subject === "string" && candidate.subject.length > 0;
  const hasText = typeof candidate.text === "string" && candidate.text.length > 0;
  const hasHtml = typeof candidate.html === "string" && candidate.html.length > 0;

  return hasTo && hasSubject && (hasText || hasHtml);
};

const startWorker = async (): Promise<void> => {
  console.log(`Email worker listening on queue: ${rabbitMqConfig.queue}`);

  await consumeQueue(rabbitMqConfig.queue, async (payload) => {
    if (!isEmailPayload(payload)) {
      throw new Error("Invalid message payload for email queue.");
    }

    await sendEmail(payload);
    console.log(`Email sent to ${payload.to} (${payload.subject})`);
  });
};

const shutdown = async (signal: string): Promise<void> => {
  console.log(`Received ${signal}, closing RabbitMQ connection...`);
  await closeRabbit();
  process.exit(0);
};

process.on("SIGINT", () => {
  void shutdown("SIGINT");
});

process.on("SIGTERM", () => {
  void shutdown("SIGTERM");
});

void startWorker().catch(async (error) => {
  console.error("Email worker failed to start", error);
  await closeRabbit();
  process.exit(1);
});
