import { HttpService } from '@nestjs/axios';
import { Injectable, Logger } from '@nestjs/common';
import axios, { AxiosError } from 'axios';
import { catchError, map } from 'rxjs';

@Injectable()
export class BlocksService {
  private readonly logger = new Logger(BlocksService.name);
  constructor(private readonly httpService: HttpService) {}

  async findLatest(): Promise<any> {
    try {
      this.logger.log(`-- FIND ALL --`);

      const res = this.httpService
        .get('http://nginx/api/blocks')
        .pipe(map((res) => {
          // console.log(res.data);
          
          // const { data } = res.data
          // res.data.data = JSON.parse(data ?? "{}")
          return res.data
        }))
        .pipe(
          catchError((error) => {
            this.logger.error(error);
            throw 'An error happened!';
          }),
        );

      return res;
    } catch (error) {
      this.logger.error(`-- ERROR FIND ALL ::: --`);
    }
  }

  async create(data: any): Promise<any> {
    try {
      this.logger.log(`-- CREATE --`);

      const res = this.httpService
        .post('http://nginx/api/blocks', { data: JSON.stringify(data) })
        .pipe(map((res) => {
          // const { data } = res.data
          // res.data.data = JSON.parse(data ?? "{}")
          return res.data
        }))
        .pipe(
          catchError((error) => {
            this.logger.error(error);
            throw 'An error happened!';
          }),
        );

      return res;
    } catch (error) {
      this.logger.error(`-- ERROR CREATE ::: --`);
    }
  }

  async findAll(): Promise<any> {
    try {
      this.logger.log(`-- LIST --`);

      const res = this.httpService
        .get('http://nginx/api/blockchains')
        .pipe(map((res) => {
          // const { data } = res.data
          // res.data.data = data.map((block) => 
          //   JSON.parse(block.data ?? "{}")
          // )
          return res.data
        }))
        .pipe(
          catchError((error) => {
            this.logger.error(error);
            throw 'An error happened!';
          }),
        );

      return res;
    } catch (error) {
      this.logger.error(`-- ERROR LIST ::: --`);
    }
  }
}
