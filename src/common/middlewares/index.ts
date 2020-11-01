import * as koaBody from 'koa-body';
import * as cors from '@koa/cors';
import * as helmet from 'koa-helmet';
// import { log } from './log';
import { errors } from './errros';

export const middlewares = [
  errors,
  koaBody({ multipart: true }),
  // log,
  helmet(),
  cors(),
];
