export class Post {
  public id: string;
  public data: {
    id: string;
    content: string;
    contentType: string;
  };
  public description: string;
  public likedBy: string[];
  public comments: string[];
  public username: string;

  constructor(data: any) {
    this.data = {
      id: data.data.name,
      content: '',
      contentType: data.data.contentType,
    };

    this.id = data.id;
    this.likedBy = data.likedBy;
    this.comments = data.comments;
    this.description = data.description;
    this.username = data.username;
  }
}
