$(document).ready(function() {

  $.fn.dataTable.Api.register('row().show()', function() {
    var page_info = this.table().page.info();
    // Get row index
    var new_row_index = this.index();
    // Row position
    var row_position = this.table().rows()[0].indexOf( new_row_index );
    // Already on right page ?
    if( row_position >= page_info.start && row_position < page_info.end ) {
        // Return row object
        return this;
    }
    // Find page number
    var page_to_display = Math.floor( row_position / this.table().page.len() );
    // Go to that page
    this.table().page( page_to_display );
    // Return row object
    return this;
  });

  $("#timeslot-list").hide();
  $("#timeslot-latest").show();

  var filterSeason = function(object) {
    $(".timeslot-filter-season").removeClass("active");
    $(object).addClass("active");
    $("#timeslot-latest").hide();
    $("#timeslot-season-title").text("Season " + $(object).find("td:nth-of-type(2)").html());
    timeslotTable
      .columns( 0 )
      .search( $(object).data("seasonid") )
      .draw();

    $("#timeslot-list").show();
  };

  //Must run before dataTables to apply to all season/timeslot pages.
  $(".timeslot-filter-season").on( "click", function () {
      filterSeason(this);
  });
  $("#timeslot-list tbody tr").on("click", function () {
      window.location.href = $(this).find("a").attr("href");
  });

  var seasonTable = $("#seasons").DataTable({
    "info": false,
    "lengthChange": false,
    "pageLength": 5,
    "order": [[1, "asc"]],
    "sDom": "<\"top\"i>rt<\"bottom\"lp><\"clear\">",
    "columnDefs": [
    {
      "targets": [ 2,3 ],
      "orderable": false
    }
    ],
    drawCallback: function(settings) {
      var pagination = $(this).closest('.dataTables_wrapper').find('.dataTables_paginate');
      pagination.toggle(this.api().page.info().pages > 1);
    }
  });
  var timeslotTable = $("#timeslots").DataTable({
    "info": false,
    "lengthChange": false,
    "sDom": '<"top"i>rt<"bottom"lp><"clear">',
    "pageLength": 5,
    "order": [[2, "asc"]],
    "columnDefs": [
    {
      "targets": [ 0 ],
      "visible": false,
      "searchable": true,
    },
    {
      "targets": [ 1,2 ],
      "searchable": false,
    },
    {
      "targets": [ 0,1,3 ],
      "orderable": false
    }
    ],
    drawCallback: function(settings) {
      var pagination = $(this).closest('.dataTables_wrapper').find('.dataTables_paginate');
      pagination.toggle(this.api().page.info().pages > 1);
    }
  });

  $.urlParam = function (name) {
    var results = new RegExp('[\?&]' + name + '=([^&#]*)').exec(window.location.href);
    if (results == null) {
      return null;
    }
    else {
      return decodeURI(results[1]) || 0;
    }
  };

  var getSeasonFilter = function() {
    if ($.urlParam("seasonID") != null) {
      var seasonTable = $("#seasons").dataTable();
      var rowId = seasonTable.fnFindCellRowIndexes($.urlParam("seasonID"),0);
      $("#seasons").DataTable().row(rowId).show().draw(false);
      filterSeason($("tr[data-seasonid='" + $.urlParam("seasonID") + "']"));
      window.location.href = window.location.href.split("#")[0] + "#timeslot-list";
    }
  }
  getSeasonFilter();
});
