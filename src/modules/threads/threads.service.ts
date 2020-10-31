import { ThreadsRepository } from './threads.repository';

export class ThreadsService {
  repository: ThreadsRepository;

  constructor(repository: ThreadsRepository) {
    this.repository = repository;
  }

  getAll() {
    return this.repository.getAll();
  }

  createOne() {
    return this.repository.createOne();
  }
}
