import { Context, Next } from 'koa';

export const log = async (ctx: Context, next: Next) => {
  console.log(ctx.method, ctx.url);
  await next();
};
