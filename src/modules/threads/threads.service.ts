import { PostData, PostsRepository } from './repositories/posts.repository';
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

  getAll(page = 1) {
    return this.threadsRepository.getAll(page);
  }

  getOne(id: number) {
    // TODO: add thread existence check
    return this.threadsRepository.getOne(id);
  }

  async createOne(data: Omit<PostData, 'threadId'>) {
    const thread = await this.threadsRepository.createOne();
    const post = await this.postsRepository.createOne({
      ...data,
      threadId: thread.id,
    });
    thread.posts = [post];
    return thread;
  }

  async createPost(threadId: number, data: Omit<PostData, 'threadId'>) {
    // TODO: add thread existence check
    await this.postsRepository.createOne({
      ...data,
      threadId,
    });
    return this.postsRepository.getThreadPosts(threadId);
  }
}
