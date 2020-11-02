import * as fs from 'fs-extra';
import * as path from 'path';
import { EntityRepository, Repository } from 'typeorm';
import { File } from '../entities/file.entity';

export type FileData = Pick<File, 'name' | 'extension' | 'postId'>;

@EntityRepository(File)
export class FilesRepository extends Repository<File> {
  static checkOrCreateUploadsDir() {
    const dirPath = path.resolve(__dirname, '../../../../uploads');
    return fs.ensureDir(dirPath);
  }

  createBulk(data: FileData[]) {
    return this.save(data);
  }
}
