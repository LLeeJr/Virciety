export class Post {
  public id: string;
  public data: {
    id: string;
    content: string;
    contentType: string;
    fileUrl: string | ArrayBuffer | null;
  };
  public description: string;
  public likedBy: string[];
  public comments: string[];

  constructor(getPost: any) {
    this.data = {
      id: getPost.data.id,
      content: getPost.data.content,
      contentType: getPost.data.contentType,
      fileUrl: ''
    };

    const blob = this.base64ImageToBlob(getPost.data.contentType, getPost.data.content)
    this.blobToFileUrl(blob)

    this.id = getPost.id;
    this.likedBy = getPost.likedBy;
    this.comments = getPost.comments;
    this.description = getPost.description;
  }

  blobToFileUrl(blob: Blob) {
    const reader = new FileReader();
    reader.readAsDataURL(blob);

    reader.onload = () => {
      this.data.fileUrl = reader.result;
    }
  }

  base64ImageToBlob(type: string, content: string): Blob {
    // decode base64
    const imageContent = atob(content);

    // create an ArrayBuffer and a view (as unsigned 8-bit)
    const buffer = new ArrayBuffer(imageContent.length);
    const view = new Uint8Array(buffer);

    // fill the view, using the decoded base64
    for(let n = 0; n < imageContent.length; n++) {
      view[n] = imageContent.charCodeAt(n);
    }

    // convert ArrayBuffer to Blob
    return new Blob([buffer], {type: type});
  }
}
