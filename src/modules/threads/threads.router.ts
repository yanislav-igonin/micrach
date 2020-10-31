import { Context } from 'koa';
import * as Router from 'koa-router';
import { db } from '../../common/db';
import { ThreadsController } from './threads.controller';
import { ThreadsRepository } from './threads.repository';
import { ThreadsService } from './threads.service';

const router = new Router({
  prefix: '/threads',
});

const getController = () => {
  const repository = db.getCustomRepository(ThreadsRepository);
  const service = new ThreadsService(repository);
  const controller = new ThreadsController(service);
  return controller;
};

router.get('/', async (ctx: Context) => {
  const controller = getController();
  await controller.getAll(ctx);
});

router.post('/', async (ctx: Context) => {
  const controller = getController();
  await controller.createOne(ctx);
});

export { router as threads };
