import { PostsRepository } from './repositories/posts.repository';
import { ThreadsRepository } from './repositories/threads.repository';

export class ThreadsService {
  threadsRepository: ThreadsRepository;
  postsRepository: PostsRepository;

  constructor(
    threadsRepository: ThreadsRepository,
    postsRepository: PostsRepository,
  ) {
    this.threadsRepository = threadsRepository;
    this.postsRepository = postsRepository;
  }

  getAll() {
    return this.threadsRepository.getAll();
  }

  createOne() {
    return this.threadsRepository.createOne();
  }
}
