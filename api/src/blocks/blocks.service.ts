import { HttpService } from '@nestjs/axios';
import { Injectable, Logger } from '@nestjs/common';

@Injectable()
export class BlocksService {
  private readonly logger = new Logger(BlocksService.name);
  constructor(private readonly httpService: HttpService) {}

  async findAll(): Promise<any> {
    try {
      this.logger.log(`-- FIND ALL --`);

      return this.httpService.get('http://localhost/api/blocks');
    } catch (error) {
      this.logger.error(`-- ERROR FIND ALL ::: ${error} --`);
    }
  }
}
