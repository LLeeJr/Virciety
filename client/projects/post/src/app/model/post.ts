export class Post {
  public id: string;
  public data: {
    id: string;
    content: string;
    contentType: string;
    fileUrl: string;
  };
  public description: string;
  public likedBy: string[];
  public comments: string[];

  constructor(getPost: any) {
    this.data = {
      id: getPost.data.id,
      content: '',
      contentType: getPost.data.contentType,
      fileUrl: '',
    };

    this.id = getPost.id;
    this.likedBy = getPost.likedBy;
    this.comments = getPost.comments;
    this.description = getPost.description;
  }
}
