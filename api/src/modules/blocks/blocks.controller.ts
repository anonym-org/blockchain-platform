import { Body, Controller, Get, Logger, Post } from '@nestjs/common';
import { CreateDto } from './blocks.contract';
import { ApiTags } from '@nestjs/swagger';
import { BlocksService } from './blocks.service';

@ApiTags('Blocks')
@Controller('blocks')
export class BlocksController {
  private readonly logger = new Logger();
  constructor(private readonly blocksService: BlocksService) {}

  @Get()
  async findOneLatestBlock(): Promise<any> {
    return this.blocksService.findLatest();
  }

  @Post()
  async createBlock(@Body() createDto: CreateDto): Promise<any> {
    return this.blocksService.create(createDto);
  }

  @Get('/history')
  async findAllBlockchains(): Promise<any> {
    return this.blocksService.findAll();
  }
}
