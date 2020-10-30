import 'reflect-metadata';
import * as Koa from 'koa';
import { config } from './common/config';
import { middlewares } from './common/middlewares';
import { routers } from './router';

const app = new Koa();

middlewares.forEach((m) => app.use(m));
routers.forEach((r) => {
  app.use(r.routes());
  app.use(r.allowedMethods());
});

app.listen(config.app.port);
console.log('app listenning on port', config.app.port);
console.log('all systems nominal');
