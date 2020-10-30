import Koa from 'koa';

import { config } from './config';

const app = new Koa();

app.listen(config.app.port);
console.log('app listenning on port', config.app.port);
console.log('all systems nominal');
