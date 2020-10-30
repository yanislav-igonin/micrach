const thread = (id: number) => ({ id, createdAt: (new Date()).toISOString() });

export class ThreadsService {
  threads = [thread(1), thread(2)];

  async getAll() {
    return this.threads;
  }
}
