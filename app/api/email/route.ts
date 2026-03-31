import { rabbitMqConfig } from "@/lib/env";
import type { EmailPayload } from "@/lib/email";
import { publishToQueue } from "@/lib/rabbitmq";

export const runtime = "nodejs";

type EnqueueResponse = {
  ok: boolean;
  queue: string;
  message: string;
};

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

export async function POST(request: Request): Promise<Response> {
  let body: unknown;

  try {
    body = await request.json();
  } catch {
    return Response.json(
      {
        ok: false,
        queue: rabbitMqConfig.queue,
        message: "Body JSON tidak valid.",
      } satisfies EnqueueResponse,
      { status: 400 },
    );
  }

  if (!isEmailPayload(body)) {
    return Response.json(
      {
        ok: false,
        queue: rabbitMqConfig.queue,
        message: "Payload wajib berisi to, subject, dan text/html.",
      } satisfies EnqueueResponse,
      { status: 422 },
    );
  }

  try {
    await publishToQueue(rabbitMqConfig.queue, body);

    return Response.json(
      {
        ok: true,
        queue: rabbitMqConfig.queue,
        message: "Email berhasil masuk ke antrian.",
      } satisfies EnqueueResponse,
      { status: 202 },
    );
  } catch (error) {
    console.error("Failed to publish email job", error);

    return Response.json(
      {
        ok: false,
        queue: rabbitMqConfig.queue,
        message: "Gagal memasukkan email ke antrian.",
      } satisfies EnqueueResponse,
      { status: 500 },
    );
  }
}
