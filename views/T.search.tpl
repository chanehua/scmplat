{{define "search"}}
<style type="text/css">
	#searchForm .row{
    			margin-bottom: 5px;
    		}
</style>
<div class="container">
	<form class="form-horizontal" id="searchForm" role="form" >
		<input type="hidden" id="p" name="p" value="1">
		<div class="row">
			<div class="col-md-4">
				<div class="row">
					<label for="sprojectName" class="control-label col-md-4">项目名:</label>
					<div class="col-md-8">
						<input class="form-control" type="text" style="width:150px;" id="sprojectName" name="sprojectName"></div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="row">
					<label for="sversion" class="control-label col-md-4">版本号:</label>
					<div class="col-md-8">
						<input class="form-control" type="text" style="width:150px;" id="sversion" name="sversion"></div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="row">
					<label for="sType" class="control-label col-md-4">类型:</label>
					<div class="col-md-8">
						<select class="form-control" style="width:150px;" id="sType" name="sType" > 
							<option></option>
							<option>stop</option>
							<option>start</option>
							<option>prc</option>
							<option>patch</option>
							<option>plat</option>
						</select> 
					</div>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-4">
				<div class="row">
					<label for="startTime" class="control-label col-md-4">开始时间:</label>
					<div class="col-md-8">
						<input class="form-control" type="text" style="width:165px;" id="startTime" name="startTime"></div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="row">
					<label for="endTime" class="control-label col-md-4">结束时间:</label>
					<div class="col-md-8">
						<input class="form-control" type="text" style="width:165px;" id="endTime" name="endTime"></div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="row">
					<label for="stargetSer" class="control-label col-md-4">目标服务器:</label>
					<div class="col-md-8">
						<input class="form-control" type="text" style="width:150px;" id="stargetSer" name="stargetSer"></div>
				</div>
			</div>
		</div>
		<button type="submit" class="btn btn-primary" style="position: absolute;
				right: 178px;" >查询</button>
	</form>
	
</div>
{{end}}