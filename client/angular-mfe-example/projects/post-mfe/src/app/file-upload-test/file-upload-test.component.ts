import {ChangeDetectorRef, Component, OnInit, OnDestroy} from '@angular/core';
import {GQLService} from "../service/g-q-l.service";

@Component({
  selector: 'app-file-upload-test',
  templateUrl: './file-upload-test.component.html',
  styleUrls: ['./file-upload-test.component.css']
})
export class FileUploadTestComponent implements OnInit, OnDestroy {
  fileBase64: any;
  file: any;
  fileBackend: any;
  description: string = '';
  content_type: string = '';

  constructor(private gqlService: GQLService,
              private cd: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.gqlService.postCreated().subscribe(
      ({ data }) => {
      console.log('got data ', data)
      },
      (error) => {
        console.log('there was an error sending the query', error)
      })
  }

  ngOnDestroy() {
    this.gqlService.closeWebSocket();
  }

  onFileSelected(event: any) {
    this.fileBackend = null;
    this.file = event.target.files[0];

    if (this.file) {
      const reader = new FileReader();
      reader.readAsDataURL(this.file);

      reader.onload = () => {
        if (reader.result) {
          const base64 = reader.result;
          const data: string = base64.toString().split(";base64,")[0];

          this.content_type = data.split(":")[1];

          this.fileBase64 = base64;
        }
      }
    }
  }

  upload() {
    if (this.fileBase64) {
      this.gqlService.upload(this.fileBase64).subscribe((data: any) => {
        console.log('got data', data.data)

        const blob = this.base64ImageToBlob(data.data.upload.contentType, data.data.upload.content)

        const reader = new FileReader();
        reader.readAsDataURL(blob);

        reader.onload = () => {
          if (reader.result) {
            const base64 = reader.result;
            const data: string = base64.toString().split(";base64,")[0];

            this.content_type = data.split(":")[1];

            this.fileBackend = base64;
          }
        }
      }, (error) => {
        console.error('there was an error sending the query', error)
      })
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

  alertFunction() {
    alert(`Content-Type ${this.content_type} is not supported`);
    this.content_type = '';
    this.cd.detectChanges();
  }
}
