import { Context, Next } from 'koa';

export const errors = async (ctx: Context, next: Next) => {
  try {
    await next();
  } catch (err) {
    console.error(err);
    ctx.body = {
      error: {
        code: 500,
        message: 'Internal Server Error',
      },
    };
  }
};
