import {
  Entity,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  UpdateDateColumn,
  OneToMany,
} from 'typeorm';
import { Post } from './post.entity';

@Entity('threads')
export class Thread {
  @PrimaryGeneratedColumn()
  id!: number;

  @OneToMany(() => Post, (post) => post.thread)
  posts!: Post[];

  @CreateDateColumn()
  createdAt!: string;

  @UpdateDateColumn()
  updatedAt!: string;
}
