import * as Router from 'koa-router';
import {
  getAll,
} from './controller';

const threads = new Router({
  prefix: '/threads',
});

threads.get('/', getAll);

export { threads };
