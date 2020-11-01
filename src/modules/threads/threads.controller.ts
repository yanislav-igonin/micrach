import { RouterContext } from 'koa-router';
import { PostData } from './repositories/posts.repository';
import { ThreadsService } from './threads.service';

export class ThreadsController {
  service: ThreadsService;

  constructor(service: ThreadsService) {
    this.service = service;
  }

  async getAll(ctx: RouterContext) {
    const { query } = ctx;
    // TODO: add query validation
    const { page } = query;
    ctx.body = await this.service.getAll(page);
  }

  async getOne(ctx: RouterContext) {
    const { id } = ctx.params;
    ctx.body = await this.service.getOne(id);
  }

  async createOne(ctx: RouterContext) {
    // TODO: add dto validation
    const dto: Omit<PostData, 'threadId'> = ctx.request.body;
    ctx.body = await this.service.createOne(dto);
  }
}
