import { Context } from 'koa';
import { PostData } from './repositories/posts.repository';
import { ThreadsService } from './threads.service';

export class ThreadsController {
  service: ThreadsService;

  constructor(service: ThreadsService) {
    this.service = service;
  }

  async getAll(ctx: Context) {
    ctx.body = await this.service.getAll();
  }

  async createOne(ctx: Context) {
    const dto: Omit<PostData, 'thread'> = ctx.request.body;
    ctx.body = await this.service.createOne(dto);
  }
}
