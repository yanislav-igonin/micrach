import { Context } from 'koa';
import * as Router from 'koa-router';
import { ThreadsController } from './threads.controller';
import { ThreadsService } from './threads.service';

const router = new Router({
  prefix: '/threads',
});

const service = new ThreadsService();
const controller = new ThreadsController(service);

router.get('/', (ctx: Context) => { controller.getAll(ctx); });

export { router as threads };
