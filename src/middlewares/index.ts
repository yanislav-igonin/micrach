import * as bodyParser from 'koa-bodyparser';
import * as cors from '@koa/cors';
import * as helmet from 'koa-helmet';
import { log } from './log';

export const middlewares = [
  bodyParser(),
  log,
  helmet(),
  cors(),
];
