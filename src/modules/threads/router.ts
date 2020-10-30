import * as Router from 'koa-router';

const threads = new Router({
  prefix: '/threads',
});

threads.get('/', (ctx) => { ctx.body = 'hello'; });

export { threads };
