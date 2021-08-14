export class Post {
  public id: string;
  public data: string;
  public description: string;
  public likedBy: string[];
  public comments: string[];

  constructor() {
    this.id = '';
    this.data = '';
    this.comments = [];
    this.likedBy = [];
    this.description = '';
  }
}
