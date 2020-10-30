import { Context } from 'koa';
import { ThreadsService } from './service';

export const getAll = async (ctx: Context) => {
  const service = new ThreadsService();
  ctx.body = await service.getAll();
};
