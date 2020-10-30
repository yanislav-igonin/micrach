import 'reflect-metadata';
import * as Koa from 'koa';
import { config } from './common/config';
import { middlewares } from './common/middlewares';
import { routers } from './router';
import { db } from './common/db';

async function main() {
  await db.connect();
  console.log('db - online');

  const app = new Koa();

  middlewares.forEach((m) => app.use(m));
  routers.forEach((r) => {
    app.use(r.routes());
    app.use(r.allowedMethods());
  });

  app.listen(config.app.port);
  console.log('app - online');
  console.log('all systems nominal');
}

main().catch((err) => {
  console.error('app - offline');
  console.error(err);
});
