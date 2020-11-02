import { isInt, isPositive } from 'class-validator';
import { RouterContext } from 'koa-router';
import { BadRequestException } from '../../common/exceptions';
import { PostData } from './repositories/posts.repository';
import { ThreadsService } from './threads.service';

export class ThreadsController {
  service: ThreadsService;

  constructor(service: ThreadsService) {
    this.service = service;
  }

  async getAll(ctx: RouterContext) {
    const { query } = ctx;
    const page = parseInt(query.page, 10);

    if (!isInt(page)) {
      throw new BadRequestException('Page must be an integer number');
    }

    if (!isPositive(page)) {
      throw new BadRequestException('Page must be a positive number');
    }

    ctx.body = await this.service.getAll(page);
  }

  async getOne(ctx: RouterContext) {
    const { params } = ctx;
    const id = parseInt(params.id, 10);

    if (!isInt(id)) {
      throw new BadRequestException('Thread ID be an integer number');
    }

    if (!isPositive(id)) {
      throw new BadRequestException('Thread ID must be a positive number');
    }

    ctx.body = await this.service.getOne(id);
  }

  async createOne(ctx: RouterContext) {
    // TODO: add dto validation
    const dto: Omit<PostData, 'threadId'> = ctx.request.body;
    // const { files } = ctx.request;
    ctx.body = await this.service.createOne(dto);
  }

  async createPost(ctx: RouterContext) {
    const { params } = ctx;
    const id = parseInt(params.id, 10);

    if (!isInt(id)) {
      throw new BadRequestException('Thread ID be an integer number');
    }

    if (!isPositive(id)) {
      throw new BadRequestException('Thread ID must be a positive number');
    }

    // TODO: add dto validation
    const dto: Omit<PostData, 'threadId'> = ctx.request.body;

    ctx.body = await this.service.createPost(id, dto);
  }
}
