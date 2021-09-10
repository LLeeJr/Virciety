import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {GQLService} from "../service/g-q-l.service";

@Component({
  selector: 'app-file-upload-test',
  templateUrl: './file-upload-test.component.html',
  styleUrls: ['./file-upload-test.component.css']
})
export class FileUploadTestComponent implements OnInit {
  fileBase64: any;
  file: any;
  fileBackend: any;
  description: string = '';
  content_type: string = '';

  constructor(private gqlService: GQLService,
              private cd: ChangeDetectorRef) {
  }

  ngOnInit(): void { }

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

  createPost() {
    if (this.fileBase64) {
      this.gqlService.createPost(this.fileBase64, this.description, 'user3').subscribe((data: any) => {
        // TODO add to shared lib list of posts
        console.log('got data', data.data)
      }, (error) => {
        console.error('there was an error sending the query', error)
      })
    }
  }

  alertFunction() {
    alert(`Content-Type ${this.content_type} is not supported`);
    this.content_type = '';
    this.cd.detectChanges();
  }
}
