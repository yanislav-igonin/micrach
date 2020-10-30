import * as bodyParser from 'koa-bodyparser';
import * as cors from '@koa/cors';
import * as helmet from 'koa-helmet';

export const middlewares = [
  bodyParser(),
  helmet(),
  cors(),
];
