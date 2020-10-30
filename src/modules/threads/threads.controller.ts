import { Context } from 'koa';
import { ThreadsService } from './threads.service';

export class ThreadsController {
  service: ThreadsService;

  constructor(service: ThreadsService) {
    this.service = service;
  }

  async getAll(ctx: Context) {
    ctx.body = await this.service.getAll();
  }
}
