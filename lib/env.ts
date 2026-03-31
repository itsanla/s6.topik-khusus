const toInt = (value: string | undefined, fallback: number): number => {
  if (!value) return fallback;
  const parsed = Number.parseInt(value, 10);
  return Number.isNaN(parsed) ? fallback : parsed;
};

export const rabbitMqConfig = {
  host: process.env.RABBITMQ_HOST ?? "127.0.0.1",
  port: toInt(process.env.RABBITMQ_PORT, 5672),
  user: process.env.RABBITMQ_USER ?? "guest",
  password: process.env.RABBITMQ_PASSWORD ?? "guest",
  vhost: process.env.RABBITMQ_VHOST ?? "/",
  queue: process.env.RABBITMQ_QUEUE ?? "email_queue",
};

export const mailConfig = {
  host: process.env.MAIL_HOST,
  port: toInt(process.env.MAIL_PORT, 465),
  secure: (process.env.MAIL_ENCRYPTION ?? "ssl") !== "tls",
  user: process.env.MAIL_USERNAME,
  password: process.env.MAIL_PASSWORD,
  fromAddress: process.env.MAIL_FROM_ADDRESS,
  fromName: process.env.MAIL_FROM_NAME ?? "Next App",
};
