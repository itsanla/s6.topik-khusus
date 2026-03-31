import nodemailer from "nodemailer";

import { mailConfig } from "./env";

const requiredMailEnv = [
  "MAIL_HOST",
  "MAIL_PORT",
  "MAIL_USERNAME",
  "MAIL_PASSWORD",
  "MAIL_FROM_ADDRESS",
] as const;

const getMissingMailEnv = (): string[] => {
  return requiredMailEnv.filter((name) => !process.env[name]);
};

const transporter = nodemailer.createTransport({
  host: mailConfig.host,
  port: mailConfig.port,
  secure: mailConfig.secure,
  auth: {
    user: mailConfig.user,
    pass: mailConfig.password,
  },
});

export type EmailPayload = {
  to: string;
  subject: string;
  text?: string;
  html?: string;
};

export const sendEmail = async (payload: EmailPayload): Promise<void> => {
  const missing = getMissingMailEnv();
  if (missing.length > 0) {
    throw new Error(`Missing mail environment variables: ${missing.join(", ")}`);
  }

  if (!payload.to || !payload.subject || (!payload.text && !payload.html)) {
    throw new Error("Invalid email payload. Require to, subject, and text/html body.");
  }

  await transporter.sendMail({
    from: `\"${mailConfig.fromName}\" <${mailConfig.fromAddress}>`,
    to: payload.to,
    subject: payload.subject,
    text: payload.text,
    html: payload.html,
  });
};
