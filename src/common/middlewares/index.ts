import * as bodyParser from 'koa-bodyparser';
import * as cors from '@koa/cors';
import * as helmet from 'koa-helmet';
import { log } from './log';
import { errors } from './errros';

export const middlewares = [
  errors,
  bodyParser(),
  // log,
  helmet(),
  cors(),
];
