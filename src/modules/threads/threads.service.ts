import { ThreadsRepository } from './threads.repository';

export class ThreadsService {
  repository: ThreadsRepository;

  constructor(repository: ThreadsRepository) {
    this.repository = repository;
  }

  async getAll() {
    return this.repository.getAll();
  }
}
