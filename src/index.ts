import * as Koa from 'koa';
import { config } from './config';
import { middlewares } from './middlewares';

const app = new Koa();

middlewares.forEach((m) => app.use(m));

app.listen(config.app.port);
console.log('app listenning on port', config.app.port);
console.log('all systems nominal');
