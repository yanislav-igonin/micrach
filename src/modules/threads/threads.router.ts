import * as Router from 'koa-router';
import { db } from '../../common/db';
import { ThreadsController } from './threads.controller';
import { ThreadsRepository } from './repositories/threads.repository';
import { ThreadsService } from './threads.service';
import { PostsRepository } from './repositories/posts.repository';

const router = new Router({
  prefix: '/threads',
});

const getController = () => {
  const threadsRepository = db.getCustomRepository(ThreadsRepository);
  const postsRepository = db.getCustomRepository(PostsRepository);
  const service = new ThreadsService(threadsRepository, postsRepository);
  const controller = new ThreadsController(service);
  return controller;
};

router.get('/', async (ctx: Router.RouterContext) => {
  const controller = getController();
  await controller.getAll(ctx);
});

router.get('/:id', async (ctx: Router.RouterContext) => {
  const controller = getController();
  await controller.getOne(ctx);
});

router.post('/', async (ctx: Router.RouterContext) => {
  const controller = getController();
  await controller.createOne(ctx);
});

router.post('/:id', async (ctx: Router.RouterContext) => {
  const controller = getController();
  await controller.createPost(ctx);
});

export { router as threads };
