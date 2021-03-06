import {Comment} from "./comment";

export class Post {
  public id: string;
  public data: {
    name: string;
    content: string;
    contentType: string;
  };
  public description: string;
  public likedBy: string[];
  public comments: Comment[];
  public username: string;
  public editMode: boolean;

  constructor(data: any) {
    this.data = {
      name: data.data.name,
      content: data.data.content,
      contentType: data.data.contentType,
    };

    this.id = data.id;
    this.likedBy = data.likedBy;
    this.comments = data.comments;
    this.description = data.description;
    this.username = data.username;
    this.editMode = false;
  }
}
