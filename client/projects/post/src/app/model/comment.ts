export class Comment {
  public id: string
  public postID: string
  public comment: string
  public createdBy: string
  public createdAt: string

  constructor(data: any) {
    this.id = data.id;
    this.postID = data.postID;
    this.comment = data.comment;
    this.createdBy = data.createdBy;
    this.createdAt = data.createdAt;
  }
}
