import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";

@Component({
  selector: 'app-file-upload-test',
  templateUrl: './file-upload-test.component.html',
  styleUrls: ['./file-upload-test.component.css']
})
export class FileUploadTestComponent implements OnInit {
  private UPLOAD_FILE = gql`
      mutation upload($file: Upload!) {
        upload(file: $file) {
          id
          name
          content
        }
      }
    `;

  private apollo: ApolloBase;
  imageURL: any;
  description: string = '';
  imageFile: any;


  constructor(private apolloProvider: Apollo) {
    this.apollo = this.apolloProvider.use('post');
  }

  ngOnInit(): void {
  }

  onFileSelected(event: any) {
    this.imageFile = event.target.files[0];

    if (this.imageFile) {
      const reader = new FileReader();
      reader.readAsDataURL(this.imageFile);

      reader.onload = (_event) => {
        this.imageURL = reader.result;
      }
    }
  }

  upload() {
    if (this.imageFile) {

      this.apollo.mutate({
        mutation: this.UPLOAD_FILE,
        variables: {
          file: this.imageFile
        }
      }).subscribe(({data}) => {
        console.log('got data', data)
      }, (error) => {
        console.error('there was an error sending the query', error)
      })
    }
  }
}
