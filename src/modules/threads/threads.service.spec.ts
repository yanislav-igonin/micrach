import { Post } from './entities/post.entity';
import { Thread } from './entities/thread.entity';
import { ThreadsService } from './threads.service';
import { ThreadsRepository } from './repositories/threads.repository';
import { PostsRepository } from './repositories/posts.repository';
import { NotFoundException } from '../../common/exceptions';

const getPostMock = (threadId: number, id: number): Post => ({
  id,
  threadId,
  title: 'MOCK TITLE',
  text: 'MOCK TEXT',
  isSage: false,
  files: [],
  createdAt: (new Date()).toISOString(),
});
const getThreadMock = (id: number): Thread => ({
  id,
  posts: [getPostMock(id, 1), getPostMock(id, 2)],
  createdAt: (new Date()).toISOString(),
  updatedAt: (new Date()).toISOString(),
});

describe('ThreadsService', () => {
  let service: ThreadsService;
  let threadsRepository: ThreadsRepository;
  let postsRepository: PostsRepository;

  beforeEach(async () => {
    threadsRepository = new ThreadsRepository();
    postsRepository = new PostsRepository();
    service = new ThreadsService(threadsRepository, postsRepository);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('get all', () => {
    it('should return 0 threads', async () => {
      jest
        .spyOn(threadsRepository, 'getAll')
        .mockResolvedValueOnce([]);

      const threads = await service.getAll();
      expect(threads).toHaveLength(0);
    });

    it('should return 2 threads', async () => {
      const existedThreads = [getThreadMock(1), getThreadMock(2)];

      jest
        .spyOn(threadsRepository, 'getAll')
        .mockResolvedValueOnce(existedThreads);

      const threads = await service.getAll();
      expect(threads).toHaveLength(2);
    });
  });

  describe('get one', () => {
    it('should throw NotFoundException for non-existent thread', async () => {
      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(undefined);

      try {
        await service.getOne(1);
        expect(true).toBeFalsy();
      } catch (err) {
        expect(err).toBeInstanceOf(NotFoundException);
        expect(err.message).toEqual('Thread not found');
      }
    });

    it('should return thread', async () => {
      const existedThread = getThreadMock(1);

      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(existedThread);

      const thread = await service.getOne(1);
      expect(thread).toStrictEqual(existedThread);
    });
  });

  describe('create one', () => {
    it('should return new thread with post', async () => {
      const newPost = getPostMock(1, 1);
      const newThread = getThreadMock(1);
      newThread.posts = [newPost];

      jest
        .spyOn(threadsRepository, 'createOne')
        .mockResolvedValueOnce(newThread);
      jest
        .spyOn(postsRepository, 'createOne')
        .mockResolvedValueOnce(newPost);

      const thread = await service.createOne(newPost);
      expect(thread).toStrictEqual(newThread);
    });
  });

  describe('create post', () => {
    it('should throw NotFoundException for non-existent thread', async () => {
      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(undefined);

      try {
        await service.getOne(1);
        expect(true).toBeFalsy();
      } catch (err) {
        expect(err).toBeInstanceOf(NotFoundException);
        expect(err.message).toEqual('Thread not found');
      }
    });

    it('should create new post and return all posts for this thread', async () => {
      const existedThread = getThreadMock(1);
      const newPost = getPostMock(existedThread.id, existedThread.posts.length);
      existedThread.posts.push(newPost);

      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(existedThread);
      jest
        .spyOn(postsRepository, 'createOne')
        .mockResolvedValueOnce(newPost);
      jest
        .spyOn(threadsRepository, 'updateThreadTime')
        .mockResolvedValueOnce(undefined);
      jest
        .spyOn(postsRepository, 'getThreadPosts')
        .mockResolvedValueOnce(existedThread.posts);

      const posts = await service.createPost(existedThread.id, newPost);
      expect(posts).toStrictEqual(existedThread.posts);
    });
  });
});
