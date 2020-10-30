import { Context, Next } from 'koa';

export const log = async (ctx: Context, next: Next) => {
  console.log(ctx);
  await next();
};
