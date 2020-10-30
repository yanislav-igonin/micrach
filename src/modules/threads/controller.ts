import { Context } from 'koa';

export const getAll = async (ctx: Context) => {
  ctx.body = 'hello threads';
};
