import {
  Entity,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  ManyToOne,
  Column,
  OneToMany,
} from 'typeorm';
import { File } from './file.entity';
import { Thread } from './thread.entity';

@Entity('posts')
export class Post {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column()
  threadId!: number;

  @ManyToOne(() => Thread, (thread) => thread.posts)
  thread?: Thread;

  @Column()
  title!: string;

  @Column()
  text!: string;

  @Column()
  isSage!: boolean;

  @CreateDateColumn()
  createdAt!: string;

  @OneToMany(() => File, (file) => file.post)
  files!: File[];
}
