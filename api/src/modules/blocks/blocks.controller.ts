import { Body, Controller, Get, Logger, Post } from '@nestjs/common';
import { CreateDto, TBlock } from './blocks.contract';
import { ApiTags } from '@nestjs/swagger';
import { BlocksService } from './blocks.service';

@ApiTags('Blocks')
@Controller('blocks')
export class BlocksController {
  private readonly logger = new Logger();
  constructor(private readonly blocksService: BlocksService) {}

  @Get()
  async findOneLatestBlock(): Promise<{
    message: string;
    data: TBlock;
  }> {
    return this.blocksService.findLatest();
  }

  @Post()
  async createBlock(@Body() createDto: CreateDto): Promise<{
    message: string;
    data: TBlock;
  }> {
    return this.blocksService.create(createDto);
  }

  @Get('/history')
  async findAllBlockchains(): Promise<{
    message: string;
    data: TBlock[];
  }> {
    return this.blocksService.findAll();
  }
}
