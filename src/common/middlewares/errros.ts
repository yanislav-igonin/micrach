import { Context, Next } from 'koa';
import { MicrachException } from '../exceptions';

export const errors = async (ctx: Context, next: Next) => {
  try {
    await next();
  } catch (err) {
    if (err instanceof MicrachException) {
      ctx.body = {
        error: {
          code: err.code,
          message: err.message,
        },
      };
    } else {
      console.error(err);
      ctx.body = {
        error: {
          code: 500,
          message: 'Internal Server Error',
        },
      };
    }
  }
};
