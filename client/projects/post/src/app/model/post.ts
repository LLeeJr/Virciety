export class Post {
  public id: string;
  public data: {
    name: string;
    content: string;
    contentType: string;
  };
  public description: string;
  public likedBy: string[];
  public comments: any[];
  public username: string;
  public editMode: boolean;
  public commentMode: boolean;

  constructor(data: any) {
    this.data = {
      name: data.data.name,
      content: '',
      contentType: data.data.contentType,
    };

    this.id = data.id;
    this.likedBy = data.likedBy;
    this.comments = data.comments;
    this.description = data.description;
    this.username = data.username;
    this.editMode = false;
    this.commentMode = false;
  }
}
