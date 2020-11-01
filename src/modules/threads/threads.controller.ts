import { RouterContext } from 'koa-router';
import { PostData } from './repositories/posts.repository';
import { ThreadsService } from './threads.service';

export class ThreadsController {
  service: ThreadsService;

  constructor(service: ThreadsService) {
    this.service = service;
  }

  async getAll(ctx: RouterContext) {
    ctx.body = await this.service.getAll();
  }

  async getOne(ctx: RouterContext) {
    const { id } = ctx.params;
    ctx.body = await this.service.getOne(id);
  }

  async createOne(ctx: RouterContext) {
    const dto: Omit<PostData, 'thread'> = ctx.request.body;
    ctx.body = await this.service.createOne(dto);
  }
}
