export class MicrachException extends Error {}

export class NotFoundException extends MicrachException {
  code = 404;
}
export class BadRequestException extends MicrachException {
  code = 400;
}
