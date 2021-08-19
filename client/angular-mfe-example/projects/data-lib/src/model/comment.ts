export class Comment {
  public id: string;
  public post_id: string;
  public description: string;
  public likedBy: string[];

  constructor() {
    this.id = '';
    this.post_id = '';
    this.description = '';
    this.likedBy = [];
  }
}
