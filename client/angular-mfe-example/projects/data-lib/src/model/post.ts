import {Comment} from "./comment"

export class Post {
  public id: string;
  public data: string;
  public description: string;
  public likedBy: string[];
  public comments: Comment[] | undefined;

  constructor() {
    this.id = '';
    this.data = '';
    this.comments = [];
    this.likedBy = [];
    this.description = '';
  }
}
